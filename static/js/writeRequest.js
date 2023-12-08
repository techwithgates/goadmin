
const submitForm = (e, tableName, objectId='') => {
    e.preventDefault()
    const currentData = new FormData(e.target)
    const formData = new FormData()
    const jsonData = {}

    currentData.forEach((val, key) => {
        if(val instanceof File){
            formData.append(key, val)
        }else{
            // only checkbox fields have id attrubutes
            const field = document.querySelector(`#${key}`)
            if(!field){
                jsonData[key] = val || null
            }
        }
    })

    // hanlde checkbox fields separately from other input fields
    const checkboxes = document.querySelectorAll(`#formId input[type="checkbox"]`)
    checkboxes.forEach(cb => jsonData[cb.name] = cb.checked)
    formData.append('jsonData', JSON.stringify(jsonData))

    if(objectId) {
        fetch(`/tables/${tableName}/old-object/${objectId}`, {
            method: 'PATCH',
            body: formData,
        })
        .then(res => {
            res.status === 200 && location.reload()
            return res.text()
        })
        .then(data => {
            if(!data.includes('success')){
                alert(data)
            }
        })
    }else{
        fetch(`/tables/${tableName}/new-object`, {
            method: 'POST',
            body: formData,
        })
        .then(res => {
            res.status === 201 && location.reload()
            return res.text()
        })
        .then(data => {
            if(!data.includes('success')){
                alert(data)
            }
        })
    }
}