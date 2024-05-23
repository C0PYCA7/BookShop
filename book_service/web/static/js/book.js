const pathname = window.location.pathname;
const pathParts = pathname.split('/');
const id = pathParts[pathParts.length - 2]; // Get the ID from the URL
const token = localStorage.getItem("Bearer")

if (id) {
    fetch('/book/' + id)
        .then(response => response.json())
        .then(data => {
            // Use the JSON data to populate the page
            const bookDiv = document.createElement('div');
            bookDiv.innerHTML = `
        <h1>${data.book.Name}</h1>
        <p>Жанр: ${data.book.Genre}</p>
        <p>Год: ${data.book.Year}</p>
        <p>Цена: ${data.book.Price}</p>
        <p>Автор: ${data.book.AuthorName} ${data.book.AuthorSurname}</p>
      `;
            document.getElementById('book-info').appendChild(bookDiv);

            // Add event listener for delete button
            const deleteButton = document.getElementById('delete');
            deleteButton.addEventListener('click', () => {
                if (confirm('Вы действительно желаете удалить эту книгу?')) {
                    fetch('/book/' + id, { method: 'DELETE' , headers: {'Authorization': 'Bearer '+ token}})
                        .then(response => {
                            if (response.ok) {
                                alert('Книга удалена');
                                window.location.href = '/';
                            } else {
                                alert('Ошибка при удалении книги');
                            }
                        })
                        .catch(error => console.error(error));
                }
            });
        })
        .catch(error => console.error(error));
} else {
    console.error('Book ID not found in URL');
}
