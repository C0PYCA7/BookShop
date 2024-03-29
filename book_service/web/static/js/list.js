document.addEventListener('DOMContentLoaded', () => {
    const addAuthorBtn = document.querySelector('#add-author-btn');
    const addBookBtn = document.querySelector('#add-book-btn');

    addAuthorBtn.addEventListener('click', () => {
        window.location.href = '/newauthor';
    });

    addBookBtn.addEventListener('click', () => {
        window.location.href = '/newbook';
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
