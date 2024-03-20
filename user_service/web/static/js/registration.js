document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#reg');

    form.addEventListener('submit', async (event) => {
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

        const response = await fetch('/registration', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: jsonData,
        });

        if (response.ok) {
            const data = await response.json();
            if (data.status === 200) {
                alert('Данные доставлены');
                window.location.href = '/login'; // перенаправление на страницу входа
            } else {
                alert('Ошибка при отправке данных');
            }
        } else {
            alert('Ошибка при отправке данных');
        }
    });
});
