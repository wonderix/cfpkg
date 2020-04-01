def init(self):
   self.simple = config_value(name="simple",type="string",description="Please enter a simple string")
   self.password = config_value(name="password",type="password",description="Please enter the password for your bank account")
   self.boolean = config_value(name="boolean",type="bool",description="Please enter a boolean")
   self.selection = config_value(name="selection",type="selection",description="Please select on option",options=["Schnitzel","Burger","Ice Cream"])  
