const fetch = require('node-fetch');

module.exports = {
    main: function (event, context) {
        console.log("Lambda is called!!!")

        var eventSourceURL = "<<Your_EventSource_URL>>"

        fetch(eventSourceURL + "/basesites").
    
        then(res => {
            if (res.ok) { // res.status >= 200 && res.status < 300
                console.log("(1/2) API call SUCCEEDED. There should be also a second call with the event content.")
                return res;
            } else {
                console.log("API call was unsuccessful: " + res.statusText)
                throw Error(res.statusText);
            }
        }).
        
        then(res => res.json()).
        
        then(json => {
            console.log(json)
            console.log("(2/2) API call SUCCEEDED. Above you should see CCv2 data output.")
        });
        return "Hello EC!";
    }
}