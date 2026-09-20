# Result V12 - n8n Protocol Connector Boundary

M50 documents that generic n8n protocol connector nodes are not provider API
catalog entries and do not trigger native parser work by themselves.

The public non-OpenAPI protocol notes now classify database/key-value, broker,
mail, file, directory, generic execution, generic HTTP, webhook, and GraphQL
node roots as priority signals only. The README points catalog curators to
those boundary notes.

Future protocol-family work remains possible, but only with stable metadata
artifacts and metadata-only parsing. Live introspection, protocol execution,
credential resolution, host selection, queue/topic/database selection, and
workflow runtime behavior remain outside `apitools`.
