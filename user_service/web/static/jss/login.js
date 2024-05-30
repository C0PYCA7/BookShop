document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#log')
    const regBtn = document.querySelector('#reg')

    regBtn.addEventListener('click', () => {
        window.location.href = '/registration'
    })

    form.addEventListener('submit', (event) => {
        event.preventDefault();

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
                window.location.href = '/';
            } else {
                alert(data.error);
            }
        }).catch((error) => {
            console.error('There was a problem with the fetch operation:', error);
        });
    })
})
