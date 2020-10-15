var submitPhone = function () {
  setTimeout(function () {
    var phoneStep = document.getElementById("phone-step");
    var codeStep = document.getElementById("code-step");
    var progress = document.getElementById("progress");
    var code = document.getElementById("code");
    phoneStep.className += " passed";
    codeStep.className = codeStep.className.replace("queued", "");
    progress.className = progress.className.replace("one-third", "two-third");
    code.focus();
  }, 500);

  var resendButton = document.getElementById("resend");
  var resendTitle = resendButton.value;
  var counter = 60;

  resendButton.value = resendTitle + " (" + counter + ")";
  var clock = setInterval(function () {
    counter--;
    if (counter == 0) {
      resendButton.disabled = false;
      resendButton.value = resendTitle;
      clearInterval(clock);
    } else {
      resendButton.value = resendTitle + " (" + counter + ")";
    }
  }, 1000);
};

var submitCode = function () {
  setTimeout(function () {
    var codeStep = document.getElementById("code-step");
    var successStep = document.getElementById("success-step");
    var progress = document.getElementById("progress");
    codeStep.className += " passed";
    successStep.className = successStep.className.replace("queued", "");
    progress.className = progress.className.replace("two-third", "completed");
  }, 500);
};
