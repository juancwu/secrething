import "@mantine/core/styles.css";
import "./App.css";
import {
	SignInButton,
	SignOutButton,
	SignedIn,
	SignedOut,
} from "@clerk/clerk-react";

function App() {
	return (
		<>
			<SignedOut>
				<SignInButton />
			</SignedOut>
			<SignedIn>
				<SignOutButton />
			</SignedIn>
			<p>Secrething</p>
		</>
	);
}

export default App;
