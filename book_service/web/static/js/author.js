const pathname = window.location.pathname;
const pathParts = pathname.split('/');
const id = pathParts[pathParts.length - 2]; // Get the ID from the URL
const token = localStorage.getItem("Bearer")
if (id) {
    fetch('/author/' + id)
        .then(response => response.json())
        .then(data => {
            // Use the JSON data to populate the page
            const authorDiv = document.createElement('div');
            authorDiv.innerHTML = `
        <h1>${data.author.name} ${data.author.surname}</h1>
        <p>${data.author.patronymic}</p>
        <p>${data.author.birthday}</p>
        <ul>
          ${data.author.BookList.map(book => `<li>${book}</li>`).join('')}
        </ul>
      `;

            const deleteButton = document.createElement('button');
            deleteButton.textContent = 'Удалить автора';
            deleteButton.addEventListener('click', () => {
                if (confirm('Вы уверены, что хотите удалить этого автора?')) {
                    fetch('/author/' + id, { method: 'DELETE' , headers: {'Authorization': 'Bearer '+ token}})
                        .then(response => {
                            if (response.ok) {
                                alert('Автор удален');
                                window.location.href = '/';
                            } else {
                                alert('Ошибка при удалении автора');
                            }
                        })
                        .catch(error => console.error(error));
                }
            });
            authorDiv.appendChild(deleteButton);

            document.body.appendChild(authorDiv);
        })
        .catch(error => console.error(error));
} else {
    console.error('Author ID not found in URL');
}
