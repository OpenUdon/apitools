# Result V9 - M43-M47 Source Coverage Roadmap

The roadmap now opens a source-first coverage sequence after M42:

- M43 expands AWS Smithy service coverage for SQS, DynamoDB, SES/SESv2, IAM,
  Cognito, Textract, Rekognition, Comprehend, Transcribe, ACM, ELB, and ELBv2.
- M44 expands Google Discovery coverage for BigQuery, Cloud Storage, Docs,
  Slides, Chat, YouTube, Tasks, Admin SDK, Analytics, Translate, Cloud Natural
  Language, Books, Business Profile, and Firestore where Discovery coverage is
  available.
- M45 expands Microsoft Graph and Azure REST OpenAPI coverage for Microsoft
  services represented in n8n.
- M46 reconciles high-value SaaS OpenAPI providers such as Slack, GitHub,
  Jira, Stripe, Notion, Twilio, Zendesk, Zoho, HubSpot, Salesforce, and Linear.
- M47 evaluates whether Dropbox Stone needs native parser ownership or whether
  the existing Stone-derived advisory overlay remains sufficient.

Each milestone has a dedicated `status-Mx.md` task ledger. The boundary remains
metadata-only: no API operation execution, credential resolution, signing,
token fetching, account/project/tenant/region selection, or n8n execution
support.
