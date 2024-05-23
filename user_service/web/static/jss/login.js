document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#log')

    form.addEventListener('submit', (event) => {
        event.preventDefault(); // отмена стандартного поведения формы

        const login = document.querySelector('#login').value
        const password = document.querySelector('#password').value

        const data = {
            login,
            password,
        }

        const jsonData = JSON.stringify(data)

        fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: jsonData,
        }).then((response) => {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error('Network response was not ok.');
            }
        }).then((data) => {
            if (data.status === 200) {
                alert('Данные доставлены');
                localStorage.setItem("Bearer", data.token)
                const tok = localStorage.getItem("Bearer")
                alert("token from local storage: " + tok)
                window.location.href = '/'; // перенаправление на главную страницу
            } else {
                alert(data.error); // вывод сообщения об ошибке
            }
        }).catch((error) => {
            console.error('There was a problem with the fetch operation:', error);
        });
    })
})
