import { useAuth } from "@/contexts/auth";
import { Stack, Tooltip, UnstyledButton } from "@mantine/core";
import {
	IconFingerprint,
	IconGauge,
	IconLogout,
	IconSettings,
	IconSwitchHorizontal,
	IconUser,
} from "@tabler/icons-react";
import { useRouter } from "@tanstack/react-router";
import { useState } from "react";
import classes from "./navbar.module.css";

interface NavbarLinkProps {
	icon: typeof IconGauge;
	label: string;
	active?: boolean;
	onClick?: () => void;
}

function NavbarLink({ icon: Icon, label, active, onClick }: NavbarLinkProps) {
	return (
		<Tooltip label={label} position="right" transitionProps={{ duration: 0 }}>
			<UnstyledButton
				onClick={onClick}
				className={classes.link}
				data-active={active || undefined}
			>
				<Icon size={20} stroke={1.5} />
			</UnstyledButton>
		</Tooltip>
	);
}

const mockdata = [
	{ icon: IconGauge, label: "Dashboard" },
	{ icon: IconUser, label: "Account" },
	{ icon: IconFingerprint, label: "Security" },
	{ icon: IconSettings, label: "Settings" },
];

export function Navbar() {
	const [active, setActive] = useState(2);
	const auth = useAuth();
	const router = useRouter();

	const links = mockdata.map((link, index) => (
		<NavbarLink
			{...link}
			key={link.label}
			active={index === active}
			onClick={() => setActive(index)}
		/>
	));

	const handleSignout = () => {
		if (auth.signout && window.confirm("Are you sure you want to sign out?")) {
			auth.signout().then(() => {
				router.invalidate();
			});
		}
	};

	return (
		<nav className={classes.navbar}>
			<div className={classes.navbarMain}>
				<Stack justify="center" gap={0}>
					{links}
				</Stack>
			</div>

			<Stack justify="center" gap={0}>
				<NavbarLink icon={IconSwitchHorizontal} label="Change account" />
				<NavbarLink icon={IconLogout} label="Logout" onClick={handleSignout} />
			</Stack>
		</nav>
	);
}
