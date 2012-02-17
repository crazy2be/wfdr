Go-OpenID
=========

About
-----

Go-OpenID is an implementation of OpenID in Golang.

For now, the implementation does not manage XRI identifier, and can only check authentication with a direct request.

Here are the specifications used:

* http://openid.net/specs/openid-authentication-2_0.html
* http://yadis.org/wiki/Yadis_1.0_%28HTML%29

Install
-------

	git clone http://github.com/fduraffourg/go-openid.git && cd go-openid && make && make install

Usage
-----

        var o = new(openid.OpenID)
        o.Identifier = "https://www.google.com/accounts/o8/id"
        o.Realm = "http://example.com"
        o.ReturnTo = "/loginCheck"
        url := o.GetUrl()

Now you have to redirect the user to the url returned. The OP will then forward the user back to you. To check the identity, do that:

        var o = new(openid.OpenID)
        o.ParseRPUrl(URL)
        grant, err := o.Verify()

grant is true if the user is authenticated, false otherwise. URL must contain the encoded content provided by the OP.

Once o.ParseRPUrl(URL) is executed, all the information provided by the OP are in the map o.Params. For instance you get the identity with:

     o.Params["openid.claimed_id"]

