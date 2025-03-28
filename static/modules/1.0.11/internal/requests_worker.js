let token = '';

onmessage = async function (event) {
    const {method, payload} = event.data;

    if (methods[method]) {
        const result = await methods[method](payload);
        postMessage({method, result});
    } else {
        postMessage({error: `method ${method} not found`});
    }
};

const methods = {
    setToken(t) {
        token = t;
    },
    async fetchWheelSpinResult(nextStepType) {
        let url = '';
        switch (nextStepType) {
            case 'rollCell':
                url = '/api/roll-cell';
                break;
            case 'rollWheelPreset':
                url = '/api/roll-wheel-preset';
                break;
            case 'rollItem':
                url = '/api/roll-item';
                break;
        }

        const res = await fetch(url, {
            method: "POST",
            headers: {
                "Authorization": token,
            },
        });

        if (!res.ok) return;
        const json = await res.json();

        return json.itemId;
    },
    async useItem(inventoryItemId) {
        const res = await fetch('/api/use-item', {
            method: "POST",
            headers: {
                "Authorization": token,
                "Content-type": 'application/json',
            },
            body: JSON.stringify({
                "itemId": inventoryItemId,
            }),
        });

        return res.ok ? inventoryItemId : null;
    },
    async dropItem(inventoryItemId) {
        const res = await fetch('/api/drop-item', {
            method: "POST",
            headers: {
                "Authorization": token,
                "Content-type": 'application/json',
            },
            body: JSON.stringify({
                "itemId": inventoryItemId,
            }),
        });

        return res.ok ? inventoryItemId : null;
    },
    async getNextStepType() {
        const res = await fetch('/api/get-next-step-type', {
            method: "GET",
            headers: {
                "Authorization": token,
            },
        });

        if (!res.ok) return;

        const json = await res.json();

        return json.nextStepType;
    },
}