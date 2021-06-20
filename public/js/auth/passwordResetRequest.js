const sendResetPasswordEmailEndpoint = '/password/email';
const errMsg = {
  empty: "The email can't be empty.",
  invalid: 'Please fill in the correct email.',
};

const validateEmailRegex = (email) => {
  const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(String(email).toLowerCase());
};

const getInputValue = () => {
  return { email: inputEmail.value.trim() };
};

const validateInput = ({ email } = values) => {
  let canSubmit = false;

  if (email.length == 0) {
    err_msg_email.innerText = errMsg.empty;
  } else if (!validateEmailRegex(email)) {
    err_msg_email.innerText = errMsg.invalid;
  } else {
    canSubmit = true;
    err_msg_email.innerText = '';
  }

  return canSubmit;
};

const clearInputBox = () => {
  inputEmail.value = '';
};

const sendResetPasswordEmailSucceed = async (resp) => {
  resp.json().then(function (data) {
    showNoticeMsg(data.msgHead, data.msgBody);
  });
  clearInputBox();
  return Promise.resolve();
};

const sendResetPasswordEmailFailed = async (resp) => {
  resp.json().then(function (data) {
    showErrMsg(data.errHead, data.errBody);
  });
};

const checkStatus = async (resp) => {
  if (resp.status >= 400) {
    return Promise.reject(resp);
  } else {
    return Promise.resolve(resp);
  }
};

const sendResetPasswordRequest = () => {
  submitBtn.classList.add('is-loading');

  const values = getInputValue();
  const ok = validateInput(values);

  if (!ok) {
    submitBtn.classList.remove('is-loading');
    return;
  }

  fetchData(sendResetPasswordEmailEndpoint, { body: values })
    .then(checkStatus)
    .then(sendResetPasswordEmailSucceed)
    .catch(sendResetPasswordEmailFailed)
    .finally((_) => {
      submitBtn.classList.remove('is-loading');
    });
};

const passwordResetRequestHandler = () => {
  err_msg_email = document.getElementById('err_msg_email');
  inputEmail = document.getElementsByName('email')[0];
  submitBtn = document.querySelector('#submitButton');

  submitBtn.addEventListener('click', sendResetPasswordRequest);
};

onDOMContentLoaded = (function () {
  passwordResetRequestHandler();
})();
