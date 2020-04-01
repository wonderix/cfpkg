def init(self):
  self.wssb = chart("https://chartmuseum.starkandwayne.com/charts/worlds-simplest-service-broker-1.3.1.tgz")
  self.wssb.ingress.enabled = False
  self.wssb.serviceBroker["class"]="myservice"
  self.wssb.serviceBroker.plan="shared"
  self.wssb.serviceBroker.tags="simple,shared"
  self.wssb.serviceBroker.baseGUID="A83ACAF3-ACE4-46BB-9AD2-D32EF1B9B813"
  self.wssb.serviceBroker.credentials='{"port": "4000", "host": "0.0.0.0"}'
