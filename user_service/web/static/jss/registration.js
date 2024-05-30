document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#reg');

    const logBtn = document.querySelector('#log')

    logBtn.addEventListener('click', () => {
        window.location.href = '/login'
    })

    function formatDate(date) {
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');

        return `${year}-${month}-${day}T00:00:00Z`;
    }

    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const name = document.querySelector('#name').value;
        const surname = document.querySelector('#surname').value;
        const patronymic = document.querySelector('#patronymic').value;
        const email = document.querySelector('#email').value;
        const login = document.querySelector('#login').value;
        const password = document.querySelector('#password').value;
        const birthdayInput = document.querySelector('#birthday');
        const birthday = formatDate(new Date(birthdayInput.value));

        const data = {
            name,
            surname,
            patronymic,
            email,
            login,
            password,
            birthday
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

            if (data.error === "Age is less than 14"){
                alert('Возраст должен быть более 14')
            }
            if (data.status === 200) {
                alert('Данные доставлены');
                window.location.href = '/login';
            } else {
                alert('Ошибка при отправке данных');
            }
        } else {
            alert('Ошибка при отправке данных');
        }
    });
});
