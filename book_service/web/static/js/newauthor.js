document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#newauthor');

    function formatDate(date) {
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');

        return `${year}-${month}-${day}T00:00:00Z`;
    }

    form.addEventListener('submit', (event) => {
        event.preventDefault();

        const name = document.querySelector('#name').value;
        const surname = document.querySelector('#surname').value;
        const patronymic = document.querySelector('#patronymic').value;
        const birthdayInput = document.querySelector('#birthday');
        const birthday = formatDate(new Date(birthdayInput.value));

        console.log(birthday)

        const data = {
            name,
            surname,
            patronymic,
            birthday,
        };

        const jsonData = JSON.stringify(data);

        console.log(jsonData)

        fetch('/newauthor', {
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
