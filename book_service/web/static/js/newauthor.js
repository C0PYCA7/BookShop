document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#newauthor');

    const token = localStorage.getItem("Bearer")
    alert("token from local storage: " + token)

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
        const birthdayInput = document.querySelector('#birthday');
        const birthday = formatDate(new Date(birthdayInput.value));

        const data = {
            name,
            surname,
            patronymic,
            birthday,
        };

        const jsonData = JSON.stringify(data);

        const response = await fetch('/newauthor', {
            method: 'POST',
            headers: {
                'Authorization': 'Bearer ' + token,
                'Content-Type': 'application/json'
            },
            body: jsonData,
        });

        if (response.ok) {
            const data = await response.json();

            console.log(data)

            if (data.status === 200) {
                alert('Данные доставлены');
                window.location.href = '/'; // перенаправление на главную страницу
            } else {
                alert('Ошибка при отправке данных');
            }
        } else {
            alert('Ошибка при отправке данных');
        }
    });
});
