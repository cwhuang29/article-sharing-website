const tagsConstructor = (e) => {
  if ((e.key === 'Enter' || e.keyCode === 13) && e.shiftKey) {
    if (tagsCount >= 5) {
      document.getElementById('err_msg_tags').innerText = errInputMsg.tagsTooMany;
      return;
    }

    var val = tagsInputBox.value.trim();
    if (val == '') {
      return;
    } else if (val.length > TAGS_BYTES_LIMIT) {
      // Note: this validation can't check Chinese words since 1 word takes 3 byte. Nonetheless, backend will take care of this
      document.getElementById('err_msg_tags').innerText = errInputMsg.tagsTooLong;
      return;
    }

    document.getElementById('err_msg_tags').innerText = '';
    tagsCount += 1;

    var newTag = `<span class="tag is-warning is-medium" name="tags" style="margin-right: 8px; margin-bottom: 5px">${val}<button class="delete is-small"></button></span>`;
    tagsList.innerHTML += newTag;
    tagsInputBox.value = '';
  }
};

const tagsDeconstructor = (e) => {
  if (e.target.tagName.toLowerCase() == 'button') {
    // c(e.currentTarget); The #tag-list element which registered this event listener's callback function
    tagsCount -= 1;
    e.target.parentNode.remove();
  }
};

const tagHandler = () => {
  tagsCount = 0;

  tagsInputBox = document.querySelector("input[name='tags']");
  tagsInputBox.addEventListener('keyup', tagsConstructor);

  tagsList = document.querySelector('#tags-list');
  tagsList.addEventListener('click', tagsDeconstructor);
};

onDOMContentLoaded = (function () {
  tagHandler();
})();
