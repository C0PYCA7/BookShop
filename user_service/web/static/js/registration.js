document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#reg');

    form.addEventListener('submit', (event) => {
        event.preventDefault();

        const name = document.querySelector('#name').value;
        const surname = document.querySelector('#surname').value;
        const patronymic = document.querySelector('#patronymic').value;
        const email = document.querySelector('#email').value;
        const login = document.querySelector('#login').value;
        const password = document.querySelector('#password').value;

        const data = {
            name,
            surname,
            patronymic,
            email,
            login,
            password,
        };

        const jsonData = JSON.stringify(data);

        fetch('/registration', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: jsonData,
        }).then((response) => {
            if (response.ok) {
                alert('Данные доставлены');
            } else {
                alert('Ошибка при отправке данных');
            }
        });
    });
});
