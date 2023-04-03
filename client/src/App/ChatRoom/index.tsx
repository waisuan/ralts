import React, { ChangeEvent, useEffect, useRef, useState } from "react"
import Button from "@mui/material/Button";

import { useAppSelector, useAppDispatch } from "../hooks"
import { chatSliceActions, getMessages, getUserProfile } from "./chatSlice"
import { Box, Container } from "@mui/system";
import { AppBar, Avatar, Card, CardContent, Grid, IconButton, Menu, MenuItem, Paper, TextField, Toolbar, Typography } from "@mui/material";
import SendIcon from "@mui/icons-material/Send";
import MenuIcon from "@mui/icons-material/Menu";
import AccountCircle from "@mui/icons-material/AccountCircle";
import { CredentialResponse, GoogleLogin, GoogleOAuthProvider } from "@react-oauth/google";

const ChatRoom: React.FC = () => {
    const messages = useAppSelector(getMessages)
    const dispatch = useAppDispatch()
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
    const [inputMessage, setInputMessage] = useState<String>("");
    const userProfile = useAppSelector(getUserProfile);
    const [authState, setAuthState] = useState<boolean>(localStorage.getItem('chat_sess_token') !== null);
    const messagesEndRef = useRef<null | HTMLDivElement>(null);

    const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    }

    const handleClose = () => {
        setAnchorEl(null);
    }

    const handleChange = (msg: String) => {
        setInputMessage(msg);
    }

    const handleSend = () => {
        dispatch(chatSliceActions.sendMessage(inputMessage));
        setInputMessage("");
    }

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }

    const login = (res: CredentialResponse) => {
        const authToken = res.credential;
        console.log(authToken);
        dispatch(chatSliceActions.initConnection(authToken));
    }

    useEffect(() => {
        // console.log("Checking session...");
        // const authToken = localStorage.getItem('chat_sess_token');
        //
        // console.log(userProfile);
        // console.log(authToken);
        // // TODO: check chat_sess_token exp
        // if (userProfile === undefined && authToken) {
        //     dispatch(chatSliceActions.initConnection(authToken));
        // }
        dispatch(chatSliceActions.initConnection(""));
    }, []);

    useEffect(() => {
        scrollToBottom();
    }, [messages]);

    const loginContainer = () => {
        return (
            <Box
                display="flex"
                justifyContent="center"
                alignItems="center"
                minHeight="100vh"
            >
                <Card 
                    variant="outlined"
                    sx={{
                        p: 10,
                        backgroundColor: "#33A0FF"
                    }}
                >
                    <CardContent>
                        <Typography 
                            gutterBottom
                            variant="h5"
                            component="div"
                            sx={{
                                textAlign: 'center',
                                color: 'white'
                            }}
                        >
                            Who are you?
                        </Typography>
                        <GoogleOAuthProvider clientId="DUMMY">
                            <GoogleLogin
                                onSuccess={login}
                                onError={() => {}}
                            />
                        </GoogleOAuthProvider>
                    </CardContent>
                </Card>
            </Box>
        );
    }

    const mainContainer = () => {
        return (
            <div>
                <AppBar position="static">
                    <Toolbar>
                        <IconButton
                            size="large"
                            edge="start"
                            color="inherit"
                            aria-label="menu"
                            sx={{ mr: 2 }}
                        >
                            <MenuIcon />
                        </IconButton>
                        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                            ralts.
                        </Typography>
                        <Typography>{userProfile?.name || "FakeUser"}</Typography>
                        <IconButton
                            size="large"
                            aria-label="account of current user"
                            aria-controls="menu-appbar"
                            aria-haspopup="true"
                            onClick={handleMenu}
                            color="inherit"
                        >
                            <Avatar 
                                src={userProfile?.picture}
                                sx={{ width: 30, height: 30 }}
                            />
                            {/* <AccountCircle /> */}
                        </IconButton>
                        <Menu
                            sx={{ mt: "40px" }}
                            id="menu-appbar"
                            anchorEl={anchorEl}
                            anchorOrigin={{
                                vertical: "top",
                                horizontal: "right",
                            }}
                            keepMounted
                            transformOrigin={{
                                vertical: "top",
                                horizontal: "right",
                            }}
                            open={Boolean(anchorEl)}
                            onClose={handleClose}
                        >
                            <MenuItem onClick={handleClose}>Logout</MenuItem>
                        </Menu>
                    </Toolbar>
                </AppBar>
                <Container maxWidth="md">
                    <Box sx={{
                        mt: 2,
                        width: "100%",
                        height: 800,
                        border: "1px solid black",
                        borderRadius: "10px",
                        maxHeight: 800,
                        overflow: "auto"
                    }}>
                        {messages.map((m) => (
                            <div key={m.createdAt}>
                                <Box
                                    sx={{
                                        p: "5px",
                                        m: "5px",
                                        border: "1px solid black",
                                        width: 250,
                                        backgroundColor: "#33F6FF",
                                        ...(m.username === "EvanSia" && {
                                            backgroundColor: "white",
                                            ml: "auto"
                                        })
                                    }}
                                >
                                    <Typography
                                        variant="caption"
                                        display="block"
                                        sx={{ color: "gray" }}
                                    >
                                        {m.username} | {m.createdAt}
                                    </Typography>
                                    <Typography variant="body2" display="block">
                                        {m.message}
                                    </Typography>
                                </Box>
                            </div>
                        ))}
                        <div ref={messagesEndRef} />
                    </Box>
                    <Box sx={{ mt: 2, mb: 2 }}>
                        <TextField 
                            sx={{
                                width: "100%",
                            }}
                            inputProps={{
                                maxLength: 250,
                            }}
                            InputProps={{
                                endAdornment: (
                                    <IconButton color="primary" onClick={handleSend}>
                                        <SendIcon />
                                    </IconButton>
                                )
                            }}
                            placeholder="Say something..."
                            onChange={(e) => {handleChange(e.target.value)}}
                            onKeyUp={(e) => {if (e.key === "Enter") {handleSend()}}}
                            value={inputMessage || ""}
                        />
                    </Box>
                </Container>
            </div>
        );
    }

    return (
        <div>
            { mainContainer() }
            {/*{userProfile ? mainContainer() : (authState ? <div/> : loginContainer())}*/}
        </div>
    );
}

export default ChatRoom