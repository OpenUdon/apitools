# Prompt V7 - AWS Smithy Converter Ownership

Add reusable AWS Smithy JSON to OpenAPI-shaped metadata conversion to
`apitools`, mirroring the Google Discovery ownership split.

V1 covers official AWS Smithy JSON service models using `restJson1`, `restXml`,
and `awsQuery`. Preserve SigV4 and AWS service metadata only as extensions for
downstream planning; do not implement signing, credential resolution, token
fetching, provider API execution, or AWS account/region selection.
