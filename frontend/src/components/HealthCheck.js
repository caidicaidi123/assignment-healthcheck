import React, {Component} from "react";
import Button from '@material-ui/core/Button/Button';
import WebsiteList from './WebsiteList';
var validUrl = require('valid-url');


class HealthCheck extends Component {
    constructor(props) {
        super(props);

        this.state = {
            apiUrl: 'http://localhost:8000/api/healthcheck',
            URL: '',
            websites: [],
            seconds: 0,
            isValidUrl: true
        };

        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleChange = this.handleChange.bind(this);
    }

    // page is auto refresh every minute to update all website status since server perform auto check every 5 mins
    tick() {
        window.location.reload();
        this.setState(prevState => ({
            seconds: prevState.seconds + 60
        }));
    }

    componentDidMount() {
        const {apiUrl} = this.state;
        fetch(apiUrl, {
            method: 'GET',
            headers: {'Content-Type': 'application/json'},
        })
            .then(response => response.json())
            .then(websites => {this.setState({websites})});

        this.interval = setInterval(() => this.tick(), 60000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    handleChange(event) {
        const {name, value} = event.target;

        this.setState({
            [name]: value
        })
    }

    handleSubmit(event) {
        event.preventDefault();
        const {apiUrl, URL} = this.state;

        if (URL && validUrl.isUri(URL)) {
            fetch(apiUrl, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({URL})
            });

            fetch(apiUrl, {
                method: 'GET',
                headers: {'Content-Type': 'application/json'},
            })
                .then(response => response.json())
                .then(websites => {this.setState({websites})});
            window.location.reload();
        } else {
            this.setState({isValidUrl: false})
        }
    }

    render() {
        const {URL, isValidUrl} = this.state;
        console.log(JSON.stringify(this.state.websites, undefined, 2));

        let websites;
        if (this.state.websites) {
            websites = this.state.websites.map(website => {
                return (
                    <WebsiteList
                        website={website}
                        key={website.URL}
                    />
                )
            });
        }

        return (
            <div className="container-fluid" >
                <div style={{fontSize: 40, marginBottom:30, marginTop:20, textAlign: 'center'}}>Website Health Check App</div>

                <form name="form" onSubmit={this.handleSubmit}>
                    <div className="form-group">
                        <input type="text" placeholder="Please enter URL" className="form-control" name="URL" value={URL} onChange={this.handleChange} required/>
                        {!isValidUrl
                        && <div className="help-block" style={{color:"red", }}>URL is not valid!</div>
                        }
                    </div>
                    <div className="form-group" style={{textAlign: 'center'}}>
                    <div>
                            <Button
                                variant="contained"
                                size="large"
                                style={{
                                    fontSize: 25,
                                }}
                                onClick={(e) => {this.handleSubmit(e)}}
                            >
                                Go
                            </Button>
                    </div>
                    </div>
                </form>

                <div style={{marginTop:20}}>
                    {websites}
                </div>
            </div>
        );
    }
}

export default HealthCheck
