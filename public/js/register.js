const errInputMsg = {
    empty: "This field can't be empty.",
    passwordInconsistency: "Password and confirm password does not match.",
    passwordTooShort: "Passwords must be at least 8 characters long.",
    emailFormatInvalid: "The email format is not correct.",
};
const registerEndpoint = "/register";

const setErrMsgBanner = (ele, msg) => {
    ele.innerText = msg;
};

const isValidEmail = (e) => {
    var filter = /^\s*[\w\-\+_]+(\.[\w\-\+_]+)*\@[\w\-\+_]+\.[\w\-\+_]+(\.[\w\-\+_]+)*\s*$/;
    return String(e).search(filter) != -1;
};

const validateInput = (vals) => {
    let isValid = true;
    const [fn, ln, pwd, pwdCfm, eml, gdr, mjr] = vals;

    if (vals.some((v) => v.length == 0)) {
        isValid = false;
    }

    if (fn.length == 0) {
        setErrMsgBanner(firstNameErr, errInputMsg.empty);
    } else {
        setErrMsgBanner(firstNameErr, "");
    }

    if (ln.length == 0) {
        setErrMsgBanner(lastNameErr, errInputMsg.empty);
    } else {
        setErrMsgBanner(lastNameErr, "");
    }

    if (pwd.length == 0) {
        setErrMsgBanner(passwordErr, errInputMsg.empty);
    } else if (pwd.length < 8) {
        setErrMsgBanner(passwordErr, errInputMsg.passwordTooShort);
    } else {
        setErrMsgBanner(passwordErr, "");
    }

    if (pwdCfm.length == 0) {
        setErrMsgBanner(passwordConfirmErr, errInputMsg.empty);
    } else if (pwd.length < 8) {
        setErrMsgBanner(passwordConfirmErr, errInputMsg.passwordTooShort);
    } else {
        setErrMsgBanner(passwordConfirmErr, "");
    }

    if (pwd.length > 7 && pwdCfm.length > 7) {
        if (pwd != pwdCfm) {
            isValid = false;
            setErrMsgBanner(passwordErr, errInputMsg.passwordInconsistency);
            setErrMsgBanner(passwordConfirmErr, errInputMsg.passwordInconsistency);
        } else {
            setErrMsgBanner(passwordErr, "");
            setErrMsgBanner(passwordConfirmErr, "");
        }
    }

    if (eml.length == 0) {
        setErrMsgBanner(emailErr, errInputMsg.empty);
    } else if (!isValidEmail(eml)) {
        isValid = false;
        setErrMsgBanner(emailErr, errInputMsg.emailFormatInvalid);
    } else {
        setErrMsgBanner(emailErr, "");
    }

    if (gdr.length == 0) {
        isValid = false;
        setErrMsgBanner(genderErr, errInputMsg.empty);
    } else {
        setErrMsgBanner(genderErr, "");
    }

    if (mjr.length == 0) {
        setErrMsgBanner(majorErr, errInputMsg.empty);
    } else {
        setErrMsgBanner(majorErr, "");
    }

    return isValid;
};

const retrieveInput = () => {
    let g = [...gender].filter((g) => g.checked).map((g) => g.value)[0];
    if (g === undefined) {
        // If user didn't click on any of them
        g = "";
    }
    return [firstName.value.trim(), lastName.value.trim(), password.value.trim(), passwordConfirm.value.trim(), email.value.trim(), g, major.value];
};

const checkStatus = async (resp) => {
    if (resp.status >= 400) {
        return Promise.reject(resp);
    } else {
        return Promise.resolve(resp);
    }
};

const creationSucceed = async (resp) => {
    window.location.href = resp.headers.get("Location");
    return Promise.resolve();
};

const creationFailed = async (resp) => {
    if (resp.status == 500) {
        resp.json().then(function (data) {
            showErrMsg(`<div><p><strong>Error !</strong></p><p>${data.err}</p></div>`);
        });
    } else if (resp.status == 400) {
        resp.json().then(function (data) {
            if (data.bindingError) {
                showErrMsg("<div><p><strong>Some severe errors occurred !</strong></p><p>Please reload the page and try again.</p></div>");
            } else {
                for (var key in data.err) {
                    document.getElementById(`err_msg_${key}`).innerText = data.err[key];
                }
            }
        });
    } else if (resp.status == 409) {
        resp.json().then((data) => showErrMsg(data.err));
    } else {
        showErrMsg("<div><p><strong>Some severe errors occurred !</strong></p><p>Please reload the page and try again.</p></div>");
    }
};

const postData = async (url = "", data = {}) => {
    return fetch(url, {
        method: "POST",
        mode: "cors",
        cache: "no-cache",
        credentials: "same-origin",
        headers: {
            "Content-Type": "application/json",
        },
        redirect: "follow",
        referrerPolicy: "no-referrer",
        body: JSON.stringify(data),
    });
};

const submitArticle = async (vals) => {
    const [fn, ln, pwd, pwdCfm, eml, gdr, mjr] = vals;
    postData(registerEndpoint, {
        first_name: fn,
        last_name: ln,
        password: pwd,
        email: eml,
        gender: gdr,
        major: mjr,
    })
        .then(checkStatus)
        .then(creationSucceed)
        .catch(creationFailed)
        .finally((_) => {
            submitBtn.classList.remove("is-loading");
        });
};

const submitNewPost = () => {
    submitBtn.classList.add("is-loading");

    let vals = retrieveInput();
    let res = validateInput(vals);

    if (!res) {
        submitBtn.classList.remove("is-loading");
    } else {
        submitArticle(vals);
    }
};

onDOMContentLoaded = (function () {
    showNoticeMsg("<div><strong>Reset password feature is coming out soon. Sorry for the inconvenience.</strong></div>");
    firstName = document.getElementsByName("first_name")[0];
    lastName = document.getElementsByName("last_name")[0];
    password = document.getElementsByName("password")[0];
    passwordConfirm = document.getElementsByName("password_confirm")[0];
    email = document.getElementsByName("email")[0];
    gender = document.getElementsByName("gender"); // Array
    major = document.getElementsByName("major")[0];

    firstNameErr = document.getElementById("err_msg_first_name");
    lastNameErr = document.getElementById("err_msg_last_name");
    passwordErr = document.getElementById("err_msg_password");
    passwordConfirmErr = document.getElementById("err_msg_password_confirm");
    emailErr = document.getElementById("err_msg_email");
    genderErr = document.getElementById("err_msg_gender");
    majorErr = document.getElementById("err_msg_major");

    submitBtn = document.querySelector("#submit_button");
    submitBtn.addEventListener("click", submitNewPost);
})();
