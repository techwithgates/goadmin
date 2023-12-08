
const openModalBtn = document.querySelectorAll(".openModalBtn") || null
const modal = document.getElementById("modal")
const closeModalBtn = document.getElementById("closeModalBtn")
const formContainer = document.getElementById("formContainer")
const modalTitle = document.getElementById("modalTitle")

openModalBtn && openModalBtn.forEach(btn => {
    btn.addEventListener('click', (e) => {
        modal.style.display = "block"
        const tableName = btn.getAttribute('data-tableName')
        const objectId = btn.getAttribute('data-objectId')
        
        if(objectId){
            // calling the API to generate an edit form template base on the table name
            fetch(`/tables/${tableName}/old-object/${objectId}`, {
                method: 'GET'
            })
            .then(res => {
                res.status !== 200 && console.log('Error!')
                return res.json()
            })
            .then(data => {
                formContainer.innerHTML += data['form']
                modalTitle.innerHTML += data['title']
            })
        }else{
            // calling the API to generate an add form template base on the table name
            fetch(`/tables/${tableName}/new-object`, {
                method: 'GET'
            })
            .then(res => {
                res.status !== 200 && console.log('Error!')
                return res.json()
            })
            .then(data => {
                formContainer.innerHTML += data['form']
                modalTitle.innerHTML += data['title']
            })
        }
    })
})

closeModalBtn.addEventListener("click", () => {
    modal.style.display = "none"
    formContainer.innerHTML = ''
    modalTitle.innerHTML = ''
})

window.addEventListener("click", (e) => {
    if(e.target === modal) {
        modal.style.display = "none"
        formContainer.innerHTML = ''
        modalTitle.innerHTML = ''
    }
})