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
