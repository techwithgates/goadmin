
const openModalBtn = document.querySelectorAll(".openModalBtn") || null
const modal = document.getElementById("modal")
const closeModalBtn = document.getElementById("closeModalBtn")
const formContainer = document.getElementById("formContainer")
const modalTitle = document.getElementById("modalTitle")

// pagination related variables
const objectDataArea = document.getElementById("objectDataArea")
const hasMoreArea = document.getElementById("loadMoreArea")
let offset = 0
let objectRows = ''

// function to handle pagination
const loadMore = (tableName, next) => {
    offset += parseInt(next) + 20

    fetch(`/tables/${tableName}?offset=${offset}`, {
        method: 'GET'
    })
    .then(res => {
        res.status !== 200 && console.log(res)
        return res.json()
    })
    .then(data => {
        data.Objects.forEach(id => {
            objectRows += `
                <div class="object-data">
                    <input class="object-checkbox" type="checkbox" onclick="selectObject(this.checked, ${id})">
                    <span class="openModalBtn" data-objectId="${id}" data-tableName="${tableName}">
                        ${tableName} - (ID: ${id})
                    </span>
                </div>
            `
        })

        // append the loaded data
        objectDataArea.innerHTML += objectRows
        // reset to empty string for next data retrieval
        objectRows = ''
        // remove load more area
        if(!data.HasMore) hasMoreArea.innerHTML = ''
    })
}

openModalBtn && openModalBtn.forEach(btn => {
    btn.addEventListener('click', (e) => {
        modal.style.display = "block"
        const tableName = btn.getAttribute('data-tableName')
        const objectId = btn.getAttribute('data-objectId')
        
        if(objectId){
            // call the API to generate an edit form template base on the table name
            fetch(`/tables/${tableName}/old-object/${objectId}`, {
                method: 'GET'
            })
            .then(res => {
                res.status !== 200 && console.log(res)
                return res.json()
            })
            .then(data => {
                formContainer.innerHTML += data['form']
                modalTitle.innerHTML += data['title']
            })
        }else{
            // call the API to generate an add form template base on the table name
            fetch(`/tables/${tableName}/new-object`, {
                method: 'GET'
            })
            .then(res => {
                res.status !== 200 && console.log(res)
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