export * from "@solo-io/proxy-runtime/proxy"; 
import { 
  RootContext, 
  Context, 
  registerRootContext, 
  FilterHeadersStatusValues, 
  stream_context, 
  send_local_response, 
  GrpcStatusValues,
} from "@solo-io/proxy-runtime";
import { JSON } from "assemblyscript-json";

class PluginContext extends RootContext {
  private mock_enable_: bool = false;
  createContext(context_id: u32): Context {
    return new HttpContext(context_id, this, this.mock_enable_);
  }

  onConfigure(configuration_size: u32): bool {
    if (configuration_size == 0) {
        return true;
    }
    super.onConfigure(configuration_size);
    let config = <JSON.Obj>(JSON.parse(this.getConfiguration()));
    let enable_or_null : JSON.Bool | null = config.getBool("mockEnable");
    if (enable_or_null != null) {
      this.mock_enable_ = enable_or_null.valueOf();
    }
    return true;
  }
}

class HttpContext extends Context {
  private mock_enable_: bool = false;
  constructor(context_id: u32, root_context: PluginContext, mock_enable: bool) {
    super(context_id, root_context);
    this.mock_enable_ = mock_enable;
  }
  onRequestHeaders(a: u32, end_of_stream: bool): FilterHeadersStatusValues {
    stream_context.headers.request.add("hello", "world");
    if (this.mock_enable_) {
      send_local_response(200, "", String.UTF8.encode("hello world"), [], GrpcStatusValues.Ok);
    }
    return FilterHeadersStatusValues.Continue;
  }
}

registerRootContext((context_id: u32) => { return new PluginContext(context_id); }, "");
