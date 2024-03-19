const container = document.getElementById("container");

const render = (data) => {
    const html = `
    <h1>${data.Name} ${data.Surname}</h1>
    <p>Отчество: ${data.Patronymic}</p>
    <p>Дата рождения: ${data.Birthday}</p>
    <h2>Список книг:</h2>
    <ul>
      ${data.BookList.map((book) => `<li>${book}</li>`).join("")}
    </ul>
    <p>Статус: ${data.Status}</p>
  `;
    container.innerHTML = html;
};

const getData = async () => {
    try {
        const response = await fetch("YOUR_SERVER_URL/authors/<id>"); // Replace with your actual URL
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        render(data);
    } catch (error) {
        console.error("Error fetching data:", error);
        // Handle error gracefully, e.g., display an error message to the user
    }
};

getData();
