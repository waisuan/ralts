import React, { useEffect, useRef, useState } from "react"

import { useAppSelector, useAppDispatch } from "../hooks"
import { chatSliceActions, getMessages, getUserProfile, fetchConnCount, getConnCount } from "./chatSlice"
import { Box, Container } from "@mui/system";
import { AppBar, Avatar, Button, Card, CardContent,IconButton, InputAdornment, Menu, MenuItem, TextField, Toolbar, Typography } from "@mui/material";
import SendIcon from "@mui/icons-material/Send";
import MenuIcon from "@mui/icons-material/Menu";
import { AccountCircle, ConstructionOutlined, Key } from "@mui/icons-material";
import { log } from "console";

const ChatRoom: React.FC = () => {
    const messages = useAppSelector(getMessages)
    const dispatch = useAppDispatch()
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [inputMessage, setInputMessage] = useState<String>("");
    const [username, setUsername] = useState<String>("");
    const [authToken, setAuthToken] = useState<String>("");
    const userProfile = useAppSelector(getUserProfile);
    const messagesEndRef = useRef<null | HTMLDivElement>(null);
    const connCount = useAppSelector(getConnCount);

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

    const login = () => {
        dispatch(chatSliceActions.initConnection({
            username: username,
            authToken: authToken
        }));
    }

    const logout = () => {
        dispatch(chatSliceActions.disconnect());
    }

    useEffect(() => {
        // TODO: Something funky going on here with keeping sessions alive upon page refresh...
        console.log("Checking session...");
        const creds = localStorage.getItem('chat_sess_token');
        if (creds) {
            dispatch(chatSliceActions.initConnection(JSON.parse(creds)));
            dispatch(fetchConnCount());
        }
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
                    }}
                >
                    <CardContent>
                        <Typography 
                            gutterBottom
                            variant="h5"
                            component="div"
                            sx={{
                                textAlign: 'center',
                                mb: '20px'
                            }}
                        >
                            Who are you?
                        </Typography>
                        <Box
                            sx={{
                                display: 'flex',
                                flexDirection: 'column',
                            }}
                            >
                            <TextField
                                required
                                id="username"
                                label="Name"
                                variant="standard"
                                margin="normal"
                                value={username}
                                onChange={(e) => { setUsername(e.target.value); }}
                                InputProps={{
                                    startAdornment: (
                                        <InputAdornment position="start">
                                            <AccountCircle />
                                        </InputAdornment>
                                    ),
                                }}
                            />
                            <TextField
                                required
                                id="auth-token"
                                label="Token"
                                variant="standard"
                                margin="normal"
                                type="password"
                                value={authToken}
                                onChange={(e) => { setAuthToken(e.target.value); }}
                                InputProps={{
                                    startAdornment: (
                                        <InputAdornment position="start">
                                            <Key />
                                        </InputAdornment>
                                    ),
                                }}
                            />
                            <Button 
                                variant="contained"
                                disabled={username === "" || authToken === ""}
                                onClick={login}
                                sx={{
                                    mt: '20px'
                                }}
                            >
                                Login
                            </Button>                      
                        </Box>
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
                        <Typography>{userProfile?.username}</Typography>
                        <IconButton
                            size="large"
                            aria-label="account of current user"
                            aria-controls="menu-appbar"
                            aria-haspopup="true"
                            onClick={handleMenu}
                            color="inherit"
                        >
                            <Avatar 
                                sx={{ width: 30, height: 30 }}
                            />
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
                            <MenuItem onClick={logout}>Logout</MenuItem>
                        </Menu>
                    </Toolbar>
                </AppBar>
                <Container maxWidth="md">
                    <Typography
                        variant="caption"
                        display="block"
                        sx={{ 
                            mt: 2,
                            color: "gray" 
                        }}
                    >
                        No. of users: {connCount}
                    </Typography>
                    <Box sx={{
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
                                        ...(m.username === userProfile?.username && {
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
            {userProfile ? mainContainer() : loginContainer()}
        </div>
    );
}

export default ChatRoom