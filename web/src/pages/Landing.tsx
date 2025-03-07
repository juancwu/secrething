import { Fragment, useState } from "react";
import {
	Box,
	Key,
	Lock,
	Share2,
	ShieldCheck,
	CornerDownRight,
	Package,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Link } from "@tanstack/react-router";

const features = [
	{
		id: "unique-structure",
		icon: <Key className="h-6 w-6 text-green-500" />,
		title: "Unique Structure",
		description:
			'Each bento has a unique name and can contain multiple "ingredients" (key-value pairs), making it perfect for organizing related credentials.',
	},
	{
		id: "secure-sharing",
		icon: <Share2 className="h-6 w-6 text-green-500" />,
		title: "Secure Sharing",
		description:
			"Bentos can be shared with other users through groups, allowing for team collaboration without compromising security.",
	},
	{
		id: "fine-grained-access",
		icon: <ShieldCheck className="h-6 w-6 text-green-500" />,
		title: "Fine-Grained Access",
		description:
			"Permissions control who can view or modify each bento, giving you complete control over your sensitive information.",
	},
	{
		id: "client-side-encryption",
		icon: <Lock className="h-6 w-6 text-green-500" />,
		title: "Client-Side Encryption",
		description:
			"All bento contents are encrypted on the client side before being sent to the server, ensuring your secrets remain private.",
	},
];

const bento_example = [
	{
		title: "AWS Access Key",
		secrets: [
			{
				key: "ACCESS_KEY_ID",
				value: "•••••••••••••••••",
			},
		],
	},
	{
		title: "Database Credentials",
		secrets: [
			{
				key: "DB_HOST",
				value: "dev-db.example.com",
			},
			{
				key: "DB_USER",
				value: "dev_user",
			},
			{
				key: "DB_PASSWORD",
				value: "•••••••••••",
			},
		],
	},
	{
		title: "API Keys",
		secrets: [
			{
				key: "STRIPE_TEST_KEY",
				value: "•••••••••••••••••",
			},
		],
	},
];

function Landing() {
	const [activeTab, setActiveTab] = useState("concept");

	return (
		<div className="min-h-screen bg-background text-foreground flex flex-col">
			<main className="pt-24 pb-16">
				<div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
					<div className="text-center mb-12">
						<h1 className="text-4xl font-extrabold tracking-tight ml-4 sm:text-5xl md:text-6xl">
							Secure Secret Management
							<p>
								Made <span className="text-primary">Simple</span>
							</p>
						</h1>
					</div>

					{/* Main Content */}
					<div className="mt-12">
						<Tabs
							defaultValue="concept"
							onValueChange={setActiveTab}
							value={activeTab}
							className="w-full"
						>
							<TabsList className="w-full grid-cols-3 max-w-2xl mx-auto">
								<TabsTrigger value="concept">The Concept</TabsTrigger>
								<TabsTrigger value="features">Key Features</TabsTrigger>
								<TabsTrigger value="example">See It In Action</TabsTrigger>
							</TabsList>

							<TabsContent value="concept" className="mt-8">
								<div className="grid gap-8 md:grid-cols-2 items-center">
									<div className="order-2 md:order-1">
										<h2 className="text-2xl font-bold mb-4">
											Inspired by Japanese Meal Boxes
										</h2>
										<p className="text-secondary-foreground mb-6">
											Just like a traditional Japanese bento box neatly
											organizes various food items into compartments, our
											digital "bento" organizes your sensitive information into
											secure, well-structured containers.
										</p>
										<div className="space-y-4">
											<div className="flex items-start">
												<div className="flex-shrink-0 mt-1">
													<CornerDownRight className="h-5 w-5 text-primary" />
												</div>
												<p className="ml-3 text-secondary-foreground">
													<span className="font-semibold">Organization:</span>{" "}
													Keep related secrets together in one bento.
												</p>
											</div>
											<div className="flex items-start">
												<div className="flex-shrink-0 mt-1">
													<CornerDownRight className="h-5 w-5 text-green-500" />
												</div>
												<p className="ml-3 text-secondary-foreground">
													<span className="font-semibold">
														Compartmentalization:
													</span>{" "}
													Different types of secrets stay separate but related.
												</p>
											</div>
											<div className="flex items-start">
												<div className="flex-shrink-0 mt-1">
													<CornerDownRight className="h-5 w-5 text-green-500" />
												</div>
												<p className="ml-3 text-secondary-foreground">
													<span className="font-semibold">Portability:</span>{" "}
													Share your bento with others securely.
												</p>
											</div>
										</div>
									</div>
									<div className="relative order-1 md:order-2 h-64 sm:h-80 md:h-96">
										<div className="absolute inset-0 bg-gradient-to-br from-green-100 to-green-200 rounded-lg shadow-md overflow-hidden">
											<div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2">
												<Box className="h-32 w-32 text-green-500 opacity-80" />
											</div>
											<div className="grid grid-cols-2 grid-rows-2 h-full w-full p-6">
												<div className="bg-white/70 m-2 rounded shadow-sm flex items-center justify-center">
													<Key className="h-8 w-8 text-gray-600" />
												</div>
												<div className="bg-white/70 m-2 rounded shadow-sm flex items-center justify-center">
													<Lock className="h-8 w-8 text-gray-600" />
												</div>
												<div className="bg-white/70 m-2 rounded shadow-sm flex items-center justify-center">
													<ShieldCheck className="h-8 w-8 text-gray-600" />
												</div>
												<div className="bg-white/70 m-2 rounded shadow-sm flex items-center justify-center">
													<Share2 className="h-8 w-8 text-gray-600" />
												</div>
											</div>
										</div>
									</div>
								</div>
							</TabsContent>

							<TabsContent value="features" className="mt-8">
								<div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
									{features.map((feature) => (
										<Card key={feature.id}>
											<CardContent>
												<div className="rounded-full bg-green-100 p-3 w-12 h-12 flex items-center justify-center mb-4">
													{feature.icon}
												</div>
												<h3 className="font-bold text-lg mb-2">
													{feature.title}
												</h3>
												<p className="text-secondary-foreground text-sm">
													{feature.description}
												</p>
											</CardContent>
										</Card>
									))}
								</div>
							</TabsContent>

							<TabsContent value="example" className="mt-8">
								<Card>
									<CardHeader>
										<CardTitle>
											Example: "Development Credentials" Bento
										</CardTitle>
									</CardHeader>
									<CardContent>
										<Card className="relative mb-8">
											<CardContent>
												<Badge className="absolute top-0 left-1/2 transform -translate-x-1/2 -translate-y-1/2 px-4 py-1 rounded-full text-sm font-medium">
													Bento Name
												</Badge>
												<p className="pt-2 text-center font-bold">
													Development Credentials
												</p>
											</CardContent>
										</Card>

										<div className="space-y-4">
											{bento_example.map((example) => (
												<Card key={example.title}>
													<CardContent>
														<div className="flex justify-between items-center mb-2">
															<div className="font-semibold">
																{example.title}
															</div>
															<div className="bg-gray-100 px-2 py-1 rounded text-xs">
																Key-Value
															</div>
														</div>
														<div className="grid grid-cols-2 gap-2">
															{example.secrets.map((secret) => (
																<Fragment key={secret.key}>
																	<div className="bg-gray-50 p-2 rounded text-sm">
																		{secret.key}
																	</div>
																	<div className="bg-gray-50 p-2 rounded text-sm">
																		{secret.value}
																	</div>
																</Fragment>
															))}
														</div>
													</CardContent>
												</Card>
											))}
										</div>

										<div className="mt-6 flex items-center text-sm text-secondary-foreground">
											<Lock className="h-4 w-4 mr-1" />
											<span>
												All data is encrypted on your device before being sent
												to the server
											</span>
										</div>
									</CardContent>
								</Card>
							</TabsContent>
						</Tabs>
					</div>

					{/* Call to Action */}
					<div className="mt-16 text-center">
						<h2 className="text-2xl font-bold mb-4">
							Ready to organize your secrets?
						</h2>
						<div className="flex flex-col sm:flex-row justify-center gap-4">
							<Button size="lg" asChild>
								<Link to="/register">Create Your First Bento</Link>
							</Button>
						</div>
					</div>
				</div>
			</main>

			<footer className="mt-auto border-t py-10">
				<div className="container mx-auto px-4">
					<div className="flex flex-col md:flex-row justify-between items-center">
						<div className="mb-4 md:mb-0">
							<span className="text-xl font-bold flex items-center gap-2">
								<Package className="h-5 w-5" /> Konbini
							</span>
							<p className="text-muted-foreground">
								Secure Secret Management Made Simple
							</p>
						</div>
						<div className="flex flex-col md:flex-row items-center gap-4">
							<a
								href="https://github.com/juancwu/konbini"
								className="text-muted-foreground hover:text-foreground transition-colors"
							>
								GitHub
							</a>
						</div>
					</div>
					<div className="mt-6 text-center text-muted-foreground">
						<p>Licensed under the MIT License</p>
					</div>
				</div>
			</footer>
		</div>
	);
}

export default Landing;
