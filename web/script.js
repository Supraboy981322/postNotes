const board = document.getElementById("board");
window.onload = loadNotes();

async function fetchJSONasArray(url) {
    try {
        const response = await fetch(url);
        
        if (!response.ok) {
            throw new Error(`http err:  ${response.status}`);
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error("err fetching data:  ", error);
        return [];
    }
}

function loadNotes() {
    (async () => {
        const notesArray = await fetchJSONasArray("/notes.json");
        if (Array.isArray(notesArray)) {
            console.log("success fetching notes");
            parseNotes(notesArray);
        } else {
            console.log("expected json data as an array for fetching notes, but is not data is not an array!");
        }
    })();
}

function parseNotes(notesArray) {
    console.log("parsing notes...");
 
    for (i = 0; i < notesArray.length; i++) {
        let noteTags = notesArray[i].tag;
        let noteCategory = notesArray[i].category;
        let noteText =  notesArray[i].text;

        let noteTagsFormatted = ""
        for (t = 0; t < noteTags.length; t++) {
            noteTagsFormatted += `<span class="tag">${noteTags[t]}</span>`;
        }

        createNote(noteCategory, noteTagsFormatted, noteText);
    }
    console.log("created notes");
}

function createNote(category, tag, content) {
    let noteContainer = document.createElement("div");
    noteContainer.setAttribute("class", "noteContainer");

    let note = document.createElement("div");
    note.setAttribute("class", "note");

    let noteTag = document.createElement("p");
    noteTag.innerHTML = tag;
    note.appendChild(noteTag);
    
    let noteCategory = document.createElement("p");
    noteCategory.setAttribute("class", "category");
    noteCategory.innerText = category;
    note.appendChild(noteCategory);
    
    let noteContent = document.createElement("div");
    noteContent.setAttribute("class", "content");
    noteContent.innerText = content;
    note.appendChild(noteContent);

    noteContainer.appendChild(note);
    board.appendChild(noteContainer);
}
