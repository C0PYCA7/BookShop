document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#update')

    const token = localStorage.getItem("Bearer")

    form.addEventListener('submit', async (event) => {
        event.preventDefault()
        const login = document.querySelector('#login').value
        const permission = document.querySelector('#permission').value

        const data = {
            login,
            permission,
        }

        const jsonData = JSON.stringify(data)

        const response = await fetch("/update", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
            body: jsonData
        })

        if (response.ok) {
            const data = await response.json()

            if (data.status === 200) {
                alert('Данные доставлены')
                window.location.href = '/'
            }else{
                alert('Ошибка')
            }
        }else{
            alert('Ошибка')
        }
    })
})