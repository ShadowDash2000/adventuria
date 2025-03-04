export default class Helper {
    static getFile(key, item, params) {
        let uri = "/api/files/" + item.collectionId + "/" + item.id + "/" + item[key];

        let i = 0;
        for (const param in params) {
            if (i === 0) uri += `?`;
            else uri += `&`;

            uri += `${param}=${params[param]}`;

            i++;
        }

        return uri;
    }
}