document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#newbook');

    function formatDate(date) {
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');

        return `${year}-${month}-${day}T00:00:00Z`;
    }

    form.addEventListener('submit', (event) => {
        event.preventDefault();

        const name = document.querySelector('#name').value;
        const author = document.querySelector('#author').value;
        const genre = document.querySelector('#genre').value;
        const birthdayInput = document.querySelector('#date');
        const date = formatDate(new Date(birthdayInput.value));
        const price = parseFloat(document.querySelector('#price').value);

        const data = {
            name,
            author,
            genre,
            date,
            price,
        };

        const jsonData = JSON.stringify(data);

        console.log(jsonData)

        fetch('/newbook', {
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
