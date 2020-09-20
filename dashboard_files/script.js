console.log("Starting TrackerJohn");

document.getElementById("new_operation").onclick = function(ev) {
    console.log("Creating New Operation")
    newOperationModal.style.display = "block";
};

var newOperationModal = document.getElementById("newOperationModal");

var closeSpan = document.getElementsByClassName("close")[0];

closeSpan.onclick = function() {
    newOperationModal.style.display = "none";
};

window.onclick = function(event) {
    if (event.target == newOperationModal) {
        newOperationModal.style.display = "none";
    }
};

function parseOperations(operations) {
    //skip if no operations
    if (operations["operations"].length == 0)
        return;
    let filler = "";
    let header = "<tr id=\"header\"><th></th><th>Date</th><th>Goal</th><th>Title</th><th>Core Factor</th><th>Secondary Factor</th><th>Duration</th></tr>"
    let tableRow = "<tr id=\"OPERATION\">\n<td></td>CONTENT</tr>";
    let tableCell = "<td>DATE</td><td>GOAL</td><td>TITLE</td><td>COREFACTOR</td><td>SECONDARYFACTOR</td><td><em>DURATION</em> minutes</td>";
    for (operation in operations["operations"]){
        let op = operations["operations"][operation];
        filler = filler + tableRow.replace("OPERATION", op.Title).replace("CONTENT", tableCell.replace("DATE", op.Date).replace("GOAL", op.Goal).replace("TITLE", op.Title).replace("COREFACTOR", op.CoreFactor).replace("SECONDARYFACTOR", op.SecondaryFactor).replace("DURATION", op.Duration));
    }
    document.getElementById("table_body").innerHTML = header + filler;
}

function getOperations(){
    let user = document.getElementById("first_name").innerHTML;    
    let xhr = new XMLHttpRequest();

    xhr.onload = function(){
        if (xhr.status != 200) {
            console.log("error getting operations");
        }else {
            console.log("got operations");
            var res = JSON.parse(xhr.response);
            parseOperations(res);
        }
    };

    xhr.open('GET', '/get_operations?user=' + user);

    xhr.send();
}

getOperations()
