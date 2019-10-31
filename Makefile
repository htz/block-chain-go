.PHONY: serve
serve:
	dev_appserver.py --support_datastore_emulator=true appengine/api/app.yaml
