document.addEventListener("DOMContentLoaded", function () {
  const simplemde = new SimpleMDE({
    element: document.getElementById("markdownEditor"),
    toolbar: ["bold", "italic", "heading", "|", "quote", "unordered-list", "ordered-list", "|", "code", "link"],
    spellChecker: false,
  });
});

function countCharacters(inputId, counterId) {
  var input = document.getElementById(inputId);
  var counter = document.getElementById(counterId);
  counter.textContent = input.value.length;
}