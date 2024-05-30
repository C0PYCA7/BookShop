document.addEventListener('DOMContentLoaded', () => {
    const addAuthorBtn = document.querySelector('#add-author-btn');
    const addBookBtn = document.querySelector('#add-book-btn');
    const updateBtn = document.querySelector('#update-permission')
    const logoutBtn = document.querySelector('#logout')
    const delBtn = document.querySelector('#del')
    const token = localStorage.getItem("Bearer")

    const sendAuthenticatedRequest = (url) => {
        fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        })
            .then(response => {
                if (response.status === 403){
                    alert("Недостаточно прав")
                }
                if (response.ok) {
                    window.location.href = url;
                } else {
                    console.error('Ошибка при отправке запроса:', response.statusText);
                }
            })
            .catch(error => {
                console.error('Ошибка при отправке запроса:', error);
            });
    };

    addAuthorBtn.addEventListener('click', () => {
        sendAuthenticatedRequest('/newauthor');
    });

    addBookBtn.addEventListener('click', () => {
        sendAuthenticatedRequest('/newbook');
    });

    updateBtn.addEventListener('click', () => {
        sendAuthenticatedRequest('/update')
    })

    delBtn.addEventListener('click', () => {
        sendAuthenticatedRequest('/delete')
    })

    logoutBtn.addEventListener('click', () => {
        localStorage.removeItem('Bearer')
        alert("Вы вышли из системы")
        window.location.href = '/login'
    })

    fetch('/list')
        .then(response => response.json())
        .then(data => {
            const bookListDiv = document.getElementById('book-list');
            const bookList = data.books;
            bookList.forEach(book => {
                const bookDiv = document.createElement('div');
                bookDiv.innerHTML = `
        <h2>${book.Name}</h2>
        <p>Author: ${book.AuthorName} ${book.AuthorSurname}</p>
        <p>Genre: ${book.Genre}</p>
        <p>Price: ${book.Price}</p>
      `;
                bookListDiv.appendChild(bookDiv);
            });
        })
        .catch(error => console.error(error));
});
