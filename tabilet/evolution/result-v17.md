# Result V17 - Gmail Raw Fnct Helper

`apitools` now exposes `helper/fnctspec` and `helper/gmailmsg`.

`gmailmsg` publishes the pure helper `gmail.render_raw`, a descriptor for
request-body-object invocation, and Go functions that render a Gmail API
`raw` MIME message string from recipient, subject, body/template, and optional
input data.

The helper remains metadata-safe and side-effect free. Runtime registration and
workflow execution stay in downstream trusted runtimes such as udon.
