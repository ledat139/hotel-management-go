const input = document.getElementById("image-input");
const preview = document.getElementById("preview");
let selectedFiles = [];

input.addEventListener("change", function (event) {
  selectedFiles = Array.from(event.target.files);
  preview.innerHTML = "";

  selectedFiles.forEach((file, index) => {
    if (!file.type.startsWith("image/")) return;

    const reader = new FileReader();
    reader.onload = function (e) {
      const wrapper = document.createElement("div");
      wrapper.className = "relative w-full aspect-square overflow-hidden";

      const img = document.createElement("img");
      img.src = e.target.result;
      img.className = "w-full h-full object-cover rounded shadow";

      const removeBtn = document.createElement("button");
      removeBtn.innerHTML = "&times;";
      removeBtn.className =
        "absolute top-1 right-1 bg-gray-700 text-white rounded-full w-6 h-6 flex items-center justify-center text-sm z-10";

      removeBtn.onclick = () => {
        selectedFiles.splice(index, 1);
        wrapper.remove();
        updateInputFiles();
      };

      wrapper.appendChild(img);
      wrapper.appendChild(removeBtn);
      preview.appendChild(wrapper);
    };
    reader.readAsDataURL(file);
  });
});

function updateInputFiles() {
  const dataTransfer = new DataTransfer();
  selectedFiles.forEach((file) => dataTransfer.items.add(file));
  input.files = dataTransfer.files;
}
