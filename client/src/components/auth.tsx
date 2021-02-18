import { FunctionalComponent, h } from "preact";
import { useState, useCallback } from "preact/hooks";

interface Props {
    page: string;
}

// firstLetterToUpper takes in a string and makes the first letter uppercase.
const firstLetterToUpper = (s: string) => {
    return s.charAt(0).toUpperCase() + s.slice(1);
};

interface RequestBody {
    username: string;
    password: string;
    master?: string; // the file encryption master password.
}

// Auth component handles login and register action with the server.
export const Auth: FunctionalComponent<Props> = (props: Props) => {
    // check that the user can only enter the predefined sites.
    if (props.page != "register" && props.page != "login") {
        return null;
    }
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [master, setMaster] = useState<string>("");

    const registerCall = useCallback(async () => {
        const requestBody: RequestBody = {
            username,
            password,
            master,
        };

        let response = await fetch(`http://localhost:8000/api/register`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json;charset=utf-8",
            },
            body: JSON.stringify(requestBody),
        });

        let res = await response.json();
        console.log(res);
    }, [username, password, master]);

    const loginCall = useCallback(async () => {
        const requestBody: RequestBody = {
            username,
            password,
        };

        let response = await fetch(`http://localhost:8000/api/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json;charset=utf-8",
            },
            body: JSON.stringify(requestBody),
        });

        let res = await response.json();
        console.log(res)
    }, [username, password]);

    const handleUsernameChange = (event: any) => {
        if (event.target.value !== "") setUsername(event.target.value);
    };

    const handlePasswordChange = (event: any) => {
        if (event.target.value !== "") setPassword(event.target.value);
    };

    const handleMasterChange = (event: any) => {
        if (event.target.value !== "") setMaster(event.target.value);
    };

    const handleAction = (event: any) => {
        event.preventDefault();

        if (props.page == "login") {
            loginCall();
        } else {
            registerCall();
        }
    };

    return (
        <div>
            <h1 style={{ marginTop: "10rem" }}>
                {firstLetterToUpper(props.page)}
            </h1>

            <form onSubmit={handleAction}>
                <div>
                    <input
                        value={username}
                        placeholder="Username"
                        onInput={handleUsernameChange}
                    />
                </div>
                <div>
                    <input
                        type="password"
                        value={password}
                        placeholder="Login password"
                        onInput={handlePasswordChange}
                    />
                </div>
                <div>
                    <input
                        type="password"
                        placeholder="Master password"
                        onInput={handleMasterChange}
                    />
                </div>

                <button type="submit">{firstLetterToUpper(props.page)}</button>
            </form>
        </div>
    );
};
