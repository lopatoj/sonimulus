<script lang="ts">
	import { Circle, Layer, Stage, type KonvaWheelEvent } from 'svelte-konva';
	let width: number;
	let height: number;
	let stageComponent: Stage | null = null;

	const scaleBy = 0.97;
	const onwheel = (e: KonvaWheelEvent) => {
		// stop default scrolling
		e.evt.preventDefault();
		if (!stageComponent) {
			return;
		}

		const stage = stageComponent.node;

		const oldScale = stage.scaleX();
		const pointer = stage.getPointerPosition()!;

		const mousePointTo = {
			x: (pointer.x - stage.x()) / oldScale,
			y: (pointer.y - stage.y()) / oldScale
		};

		// how to scale? Zoom in? Or zoom out?
		let direction = e.evt.deltaY > 0 ? 1 : -1;

		// when we zoom on trackpad, e.evt.ctrlKey is true
		// in that case lets revert direction
		if (e.evt.ctrlKey) {
			direction = -direction;
		}

		const newScale = direction > 0 ? oldScale * scaleBy : oldScale / scaleBy;

		stage.scale({ x: newScale, y: newScale });

		const newPos = {
			x: pointer.x - mousePointTo.x * newScale,
			y: pointer.y - mousePointTo.y * newScale
		};
		stage.position(newPos);
	};
</script>

<div bind:offsetWidth={width} bind:offsetHeight={height} class="h-full w-full bg-[#181818]">
	<Stage draggable={true} {width} {height} {onwheel} bind:this={stageComponent}>
		<Layer>
			<Circle x={width / 2} y={height / 2} radius={50} fill="blue" />
		</Layer>
	</Stage>
</div>
