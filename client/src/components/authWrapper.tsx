import { FunctionalComponent, h } from "preact";
import { useRef, useEffect, useCallback, useState } from "preact/hooks";

interface Props {
    child: FunctionalComponent;
}
function useDidMount() {
    const didMountRef = useRef(true);

    useEffect(() => {
        didMountRef.current = false;
    }, []);
    return didMountRef.current;
}

// AuthWrapper
export const AuthWrapper: FunctionalComponent<Props> = (props: Props) => {
    const didMount = useDidMount();
    const [authenticated, setAuthenticated] = useState<Boolean | null>(false);

    const fetchData = useCallback(async () => {
        const res = await fetch("http://localhost:8080/api/me", {
            method: "GET",
        });

        setAuthenticated(res.ok);
    }, []);
    useEffect(() => {
        if (didMount) {
            fetchData();
        }
    }, [didMount]);

    return (
        <div>
            {authenticated === null ? (
                <p>loading...</p>
            ) : authenticated ? (
                <div>{props.child}</div>
            ) : (
                <p>forbidden</p>
            )}
        </div>
    );
};
