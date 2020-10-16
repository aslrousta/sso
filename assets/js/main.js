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
  var clock = document.getElementById("clock");
  var clockArc = clock.getElementsByTagName("path")[0];
  var counter = 60;

  var clockTimer = setInterval(function () {
    counter--;
    if (counter == 0) {
      resendButton.disabled = false;
      resendButton.removeChild(clock);
      clearInterval(clockTimer);
    } else {
      var theta = ((60 - counter) * Math.PI) / 30;
      var x = 9 * Math.sin(theta);
      var y = -9 * Math.cos(theta);
      var largeArc = counter > 30 ? 0 : 1;
      clockArc.setAttribute(
        "d",
        "M0,0 v-9 A9,9 0 " + largeArc + ",1 " + x + "," + y + " Z"
      );
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
