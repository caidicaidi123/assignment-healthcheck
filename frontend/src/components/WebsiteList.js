import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import InboxIcon from '@material-ui/icons/Inbox';
import WebIcon from '@material-ui/icons/Web';
import DeleteIcon from '@material-ui/icons/Delete';
import DoneIcon from '@material-ui/icons/Done';
import ClearIcon from '@material-ui/icons/Clear';
import green from '@material-ui/core/colors/green';
import SvgIcon from '@material-ui/core/SvgIcon';

const styles = theme => ({
    root: {
        paddingTop: 5,
    },

});


function WebsiteList(props) {
    const { website } = props;
    if(!website) {
        return (
            <div></div>
        )
    }

    return (
            <ListItem button>
                <ListItemIcon>
                    <WebIcon />
                </ListItemIcon>
                <ListItemText primary={website.URL} />

                <div> status:
                {
                    website.Status ? (
                        <DoneIcon style={{marginRight:16, marginLeft:4}} color="primary"/>
                    ) : (
                        <ClearIcon style={{marginRight:16, marginLeft:4}} color="error"/>
                    )
                }
                </div>
                <DeleteIcon color="secondary" onClick={() => {
                    console.log(website);
                    fetch('http://localhost:8000/api/healthcheck', {
                        method: 'DELETE',
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify({URL:website.URL})
                    });

                    window.location.reload();
                }
                }/>
            </ListItem>
    );
}

WebsiteList.propTypes = {
    classes: PropTypes.object.isRequired,
}

export default withStyles(styles)(WebsiteList);
