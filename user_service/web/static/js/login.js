document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#log')

    form.addEventListener('submit', (event) => {
        event.preventDefault();

        const login = document.querySelector('#login').value
        const password = document.querySelector('#password').value

        const data = {
            login,
            password,
        }

        const jsonData = JSON.stringify(data)

        console.log(jsonData)

        fetch('/login', {
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
        })
    })
})