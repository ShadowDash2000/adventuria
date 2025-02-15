import {app} from "../app.js";

document.addEventListener('DOMContentLoaded', () => {
    const authForm = document.getElementById('auth');
    authForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = new FormData(authForm);

        const authResult = await app.pb.collection('users').authWithPassword(
            formData.get('login'),
            formData.get('password'),
        );

        if (authResult.token) {
            app.modal.close();
            app.isAuthorized = true;
            app.auth = authResult;
            await app.updateInnerField();
        }
    });
});