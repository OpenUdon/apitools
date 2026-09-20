# Status M60 - Workato-Signaled AWS Smithy Expansion

| Item | State | Notes |
|---|---|---|
| Batch source freeze | `[+]` | Used Workato-visible AWS demand only as a prioritization signal. Implementation metadata is sourced from official AWS Smithy model paths and AWS API reference docs only; no connector metadata was copied. |
| Existing AWS coverage reconciliation | `[+]` | Reconciled existing AWS Smithy rows before adding M60 IDs. New rows use distinct provider/spec IDs and reuse the native Smithy provider/candidate/security patterns. |
| Official Smithy model review | `[+]` | Recorded official AWS `api-models-aws` model paths, blob SHAs, model sizes, docs refs, source notes, and service-shape evidence in provider source refs and artifact registry rows. |
| EC2 cataloged | `[+]` | Added `aws-ec2` with official model `models/ec2/service/2016-11-15/ec2-2016-11-15.json`, service shape `com.amazonaws.ec2#AmazonEC2`, and EC2 API docs ref. |
| CloudWatch cataloged | `[+]` | Added `aws-cloudwatch` with official model `models/cloudwatch/service/2010-08-01/cloudwatch-2010-08-01.json`, service shape `com.amazonaws.cloudwatch#GraniteServiceVersion20100801`, and CloudWatch API docs ref. |
| Kinesis cataloged | `[+]` | Added `aws-kinesis` with official model `models/kinesis/service/2013-12-02/kinesis-2013-12-02.json`, service shape `com.amazonaws.kinesis#Kinesis_20131202`, and Kinesis API docs ref. |
| RDS cataloged | `[+]` | Added `aws-rds` with official model `models/rds/service/2014-10-31/rds-2014-10-31.json`, service shape `com.amazonaws.rds#AmazonRDSv19`, and RDS API docs ref. |
| SageMaker cataloged | `[+]` | Added `aws-sagemaker` with official model `models/sagemaker/service/2017-07-24/sagemaker-2017-07-24.json`, service shape `com.amazonaws.sagemaker#SageMaker`, and SageMaker API docs ref. |
| API Gateway cataloged | `[+]` | Added `aws-api-gateway` with official model `models/api-gateway/service/2015-07-09/api-gateway-2015-07-09.json`, service shape `com.amazonaws.apigateway#BackplaneControlService`, and API Gateway API docs ref. |
| Athena cataloged | `[+]` | Added `aws-athena` with official model `models/athena/service/2017-05-18/athena-2017-05-18.json`, service shape `com.amazonaws.athena#AmazonAthena`, and Athena API docs ref. |
| Bedrock cataloged | `[+]` | Added `aws-bedrock` with official model `models/bedrock/service/2023-04-20/bedrock-2023-04-20.json`, service shape `com.amazonaws.bedrock#AmazonBedrockControlPlaneService`, and Bedrock API docs ref. |
| Glue cataloged | `[+]` | Added `aws-glue` with official model `models/glue/service/2017-03-31/glue-2017-03-31.json`, service shape `com.amazonaws.glue#AWSGlue`, and Glue API docs ref. |
| GuardDuty cataloged | `[+]` | Added `aws-guardduty` with official model `models/guardduty/service/2017-11-28/guardduty-2017-11-28.json`, service shape `com.amazonaws.guardduty#GuardDutyAPIService`, and GuardDuty API docs ref. |
| KMS cataloged | `[+]` | Added `aws-kms` with official model `models/kms/service/2014-11-01/kms-2014-11-01.json`, service shape `com.amazonaws.kms#TrentService`, and KMS API docs ref. |
| Secrets Manager cataloged | `[+]` | Added `aws-secrets-manager` with official model `models/secrets-manager/service/2017-10-17/secrets-manager-2017-10-17.json`, service shape `com.amazonaws.secretsmanager#secretsmanager`, and Secrets Manager API docs ref. |
| Security Hub cataloged | `[+]` | Added `aws-securityhub` with official model `models/securityhub/service/2018-10-26/securityhub-2018-10-26.json`, service shape `com.amazonaws.securityhub#SecurityHubAPIService`, and Security Hub API docs ref. |
| Inspector2 cataloged | `[+]` | Added `aws-inspector2` with official model `models/inspector2/service/2020-06-08/inspector2-2020-06-08.json`, service shape `com.amazonaws.inspector2#Inspector2`, and Inspector2 API docs ref. |
| Review artifacts/native parser coverage | `[+]` | Added M60 Smithy rows to the artifact registry and parser fixture coverage list. Ignored local review artifacts parse as Smithy JSON when present; no Smithy-to-OpenAPI conversion was added. |
| Security metadata | `[+]` | Added SigV4 advisory overlays and security classifications for all M60 providers. The overlays describe auth metadata only and do not add credential resolution, signing, account, or region behavior. |
| Catalog/reporting verification | `[+]` | Verified M60 rows through catalog stats, security-report, security-audit, catalog check, and representative advisory output. Provider totals now include 316 providers and 29 native Smithy providers. |
| Verification | `[+]` | Ran Go tests/vet, `GOWORK=off` tests/vet, diff checks, and downstream `../openudon` and `../udon` tests before commit. |

## Source Notes

- M60 uses Workato and similar connector lists only to prioritize which AWS
  services to review.
- Official AWS Smithy JSON models are the source of truth for promoted service
  metadata.
- Smithy remains a native source family in `apitools`; this milestone must not
  reintroduce Smithy-to-OpenAPI conversion as a runtime contract.
- SigV4, AWS protocol traits, and endpoint behavior are metadata for
  downstream consumers, not executable behavior in `apitools`.

## Boundary Checks

- Do not call AWS APIs, read AWS credential stores, fetch tokens, sign
  requests, or choose AWS accounts, regions, partitions, profiles, endpoints,
  or environments.
- Do not claim OpenAPI import readiness for native Smithy services.
- Do not copy Workato, Tray, n8n, or other third-party connector metadata,
  generated code, credential behavior, examples, descriptions, icons, or
  workflow behavior.
- Treat official Smithy artifacts as untrusted metadata and validate parser
  behavior with local fixtures only.
