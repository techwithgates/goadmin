
const cbList = []
const deleteBtn = document.getElementById('deleteBtn')
const checkboxAll = document.querySelector('#checkboxArea input')
const checkBoxes = document.querySelectorAll('.object-checkbox')
const objects = document.querySelectorAll('.object-data span')

const selectObject = (checked, id) => {
    checked ? cbList.push(id) : cbList.splice(cbList.indexOf(id), 1)

    if(cbList.length > 0){
        deleteBtn.style.display = 'flex'
    }else{
        deleteBtn.style.display = 'none'
       
        checkboxAll.checked = false
    }
}

const selectAllObjects = (checked) => {
    cbList.length = 0
    checkBoxes.forEach(cb => cb.checked = checked)

    if(checked){
        objects.forEach(obj => cbList.push(obj.getAttribute('data-objectId')))
    }

    if(cbList.length > 0){
        deleteBtn.style.display = 'flex'
    }else{
        deleteBtn.style.display = 'none'
    }
}

deleteBtn && deleteBtn.addEventListener('click', e => {
    const modelName = deleteBtn.getAttribute('data-tableName')
    fetch(`/tables/${modelName}/old-object`, {
        method: 'DELETE',
        body: JSON.stringify(cbList),
    })
    .then(res => {
        if(res.status === 200){
            location.reload()
            return res.text()
        }else{
            return res.text()
        }
    })
    .then(data => {
        if(!data.includes('success')){
            alert(data)
        }
    })
})