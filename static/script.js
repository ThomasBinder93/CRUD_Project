
    const API = "/api/items";

    async function loadItems() {
        const res = await fetch(API);
        const json = await res.json();
        const items = Array.isArray(json) ? json : json?.data ?? [];

        const list = get("list");
        list.innerHTML = "<tr><th>Name</th><th>Description</th><th>Status</th><th>Actions</th></tr>";

        items.forEach(item => {
            const tr = createElement("tr");

            tr.innerHTML = `<td><input id="name-${item.id}" value="${item.name}"></td>
                            <td><input id="description-${item.id}" value="${item.description || ''}"></td>
                            <td><select id="status-${item.id}">
                                <option value="active" ${item.status === 'active' ? 'selected' : ''}>active</option>
                                <option value="completed" ${item.status === 'completed' ? 'selected' : ''}>completed</option>
                                <option value="archived" ${item.status === 'archived' ? 'selected' : ''}>archived</option>
                            </select></td>
                            <td><button onclick="updateItem(${item.id}, getValue('name-${item.id}'), getValue('description-${item.id}'), getValue('status-${item.id}'))">Ändern</button>
                                <button onclick="deleteItem(${item.id})">Löschen</button></td>`;
            list.appendChild(tr);
        });
    }

    async function createItem() {

        const name = getValue("nameInput");
        const description = getValue("descriptionInput");
        const status = getValue("statusInput");

        const res = await fetch(API, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({name, description, status})
        });

        if (!res.ok) {
            const error = await res.json();
            showMessage(`Fehler: ${error.error.details}`, true);
            return;
        }

        get("nameInput").value = "";
        get("descriptionInput").value = "";
        get("statusInput").value = "active";
        loadItems();
        showMessage("Item erfolgreich erstellt");
    }

    async function updateItem(id, name, description, status) {
        console.log("Id ", id, " Name ", name, " Description ", description, " Status ", status);

        const res = await fetch(`${API}/${id}`, {
            method: "PUT",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({name, description, status})
        });

        if (!res.ok) {
            const error = await res.json();
            showMessage(`Fehler: ${error.error.details}`, true);
            return;
        }


        loadItems();
        showMessage("Item erfolgreich aktualisiert");
    }

    async function deleteItem(id) {
        const res = await fetch(`${API}/${id}`, {
            method: "DELETE"
        });

        if (!res.ok) {
            const error = await res.json();
            showMessage(`Fehler: ${error.error.details}`, true);
            return;
        }

        loadItems();
        showMessage("Item erfolgreich gelöscht");
    }

    loadItems();

