import { IconMoon, IconSun } from "@tabler/icons-react";
import {
	ActionIcon,
	Box,
	Button,
	Flex,
	Group,
	Modal,
	useMantineColorScheme,
} from "@mantine/core";
import classes from "@/components/ui/header.module.css";
import { Link, useRouter } from "@tanstack/react-router";
import { useAuth } from "@/contexts/auth";
import { notifications } from "@mantine/notifications";
import { useState } from "react";

export function Header() {
	const [openSignOutModal, setOpenSignOutModel] = useState(false);
	const [isLoading, setIsLoading] = useState(false);
	const { toggleColorScheme, colorScheme } = useMantineColorScheme();
	const auth = useAuth();
	const router = useRouter();

	const handleSignOut = async () => {
		if (auth.signout) {
			try {
				setIsLoading(true);
				await auth.signout();
				setOpenSignOutModel(false);
				await router.invalidate();
			} catch (error) {
				console.error(error);
				notifications.show({
					title: "Oops, something went wrong",
					message: "Failed to completely sign out. Please try again.",
				});
			} finally {
				setIsLoading(false);
			}
		}
	};

	return (
		<Box pb={60}>
			<header className={classes.header}>
				<Group justify="space-between" h="100%">
					<div>
						<ActionIcon
							onClick={toggleColorScheme}
							variant="default"
							aria-label="Toggle color scheme"
							size="lg"
						>
							{colorScheme === "dark" ? (
								<IconSun style={{ width: "70%", height: "70%" }} stroke={1.5} />
							) : (
								<IconMoon
									style={{ width: "70%", height: "70%" }}
									stroke={1.5}
								/>
							)}
						</ActionIcon>
					</div>

					<Group>
						{!auth.isAuthenticated && (
							<>
								<Button component={Link} to="/signin" variant="default">
									Sign In
								</Button>
								<Button component={Link} to="/signup">
									Sign Up
								</Button>
							</>
						)}
						{auth.isAuthenticated && (
							<>
								<Button
									onClick={() => setOpenSignOutModel(true)}
									variant="default"
								>
									Sign Out
								</Button>
							</>
						)}
					</Group>
				</Group>
			</header>

			<Modal
				opened={openSignOutModal}
				onClose={() => setOpenSignOutModel(false)}
				title="Are you sure you want to sign out?"
			>
				<Flex justify="end" gap={16}>
					<Button
						loading={isLoading}
						onClick={() => setOpenSignOutModel(false)}
					>
						No, I want to stay
					</Button>
					<Button variant="default" loading={isLoading} onClick={handleSignOut}>
						Yes, sign me out
					</Button>
				</Flex>
			</Modal>
		</Box>
	);
}
