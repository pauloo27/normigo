let state = 1;

const caption = document.getElementById("caption");
const button = document.getElementById("button");

function ocr(imageURL) {
  state = 0;
  caption.disabled = true;
  button.innerText = "..."
  
  fetch(`/ocr?url=${imageURL}`, {method: "post"}).then(res => res.text()).then(text => {
    caption.value = text;
    state = 2;
    button.innerText = "Translate";
    caption.disabled = false;
  });
}

function translate() {
  state = 0;
  caption.disabled = true;
  button.innerText = "..."

  fetch(`/tr?text=${caption.value}`, {method: "post"}).then(res => res.text()).then(text => {
    caption.value = text;
    state = 2;
    button.innerText = "Translate";
    caption.disabled = false;
  });
}

function handleClick(imageURL) {
  if (state === 0) return;
  if (state === 1) {
    ocr(imageURL);
  } else {
    translate();
  }
} 
