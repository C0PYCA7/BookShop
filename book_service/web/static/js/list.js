document.addEventListener('DOMContentLoaded', () => {
    const addAuthorBtn = document.querySelector('#add-author-btn');
    const addBookBtn = document.querySelector('#add-book-btn');

    // Получение токена из localStorage
    const token = localStorage.getItem("Bearer")
    alert("token from local storage: " + token)
    // Функция для отправки запроса на защищенный роут с токеном в заголовке
    const sendAuthenticatedRequest = (url) => {
        fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        })
            .then(response => {
                if (response.ok) {
                    // Если запрос успешен, перенаправляем пользователя на указанный URL
                    window.location.href = url;
                } else {
                    // В случае ошибки обрабатываем ее
                    console.error('Ошибка при отправке запроса:', response.statusText);
                }
            })
            .catch(error => {
                console.error('Ошибка при отправке запроса:', error);
            });
    };

    addAuthorBtn.addEventListener('click', () => {
        // Отправка запроса на страницу создания нового автора с токеном в заголовке
        sendAuthenticatedRequest('/newauthor');
    });

    addBookBtn.addEventListener('click', () => {
        // Отправка запроса на страницу создания новой книги с токеном в заголовке
        sendAuthenticatedRequest('/newbook');
    });

    fetch('/list')
        .then(response => response.json())
        .then(data => {
            const bookListDiv = document.getElementById('book-list');
            const bookList = data.books;
            bookList.forEach(book => {
                const bookDiv = document.createElement('div');
                bookDiv.innerHTML = `
        <h2>${book.Name}</h2>
        <p>Genre: ${book.Genre}</p>
        <p>Price: ${book.Price}</p>
      `;
                bookListDiv.appendChild(bookDiv);
            });
        })
        .catch(error => console.error(error));
});
